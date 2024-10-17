package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/Tony36051/go-file-agent/generated/go/file_transfer/api/v1"
)

type fileServer struct {
	pb.UnimplementedFileServiceServer
}

func (s *fileServer) RequestFile(req *pb.FilePathRequest, stream pb.FileService_RequestFileServer) error {
	filePath := req.GetFilePath()
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	chunkNumber := 0

	for {
		n, err := file.Read(buf)
		if n == 0 {
			if err == io.EOF {
				break
			}
			return err
		}

		chunk := &pb.FileChunk{
			Data:        buf[:n],
			ChunkNumber: int32(chunkNumber),
			FileName:    req.GetFilePath(),
		}

		if err := stream.Send(chunk); err != nil {
			return err
		}

		chunkNumber++
	}

	return nil
}

func main() {
	port := flag.Int("port", 6565, "the server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &fileServer{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
