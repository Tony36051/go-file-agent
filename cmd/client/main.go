package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/Tony36051/go-file-agent/generated/go/file_transfer/api/v1"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:6565", "the server address")
	filePath := flag.String("file", "path/to/your/file", "the file path to request")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	request := &pb.FilePathRequest{FilePath: *filePath}

	stream, err := client.RequestFile(context.Background(), request)
	if err != nil {
		log.Fatalf("could not request file: %v", err)
	}

	outputFile, err := os.Create("output_" + *filePath)
	if err != nil {
		log.Fatalf("could not create output file: %v", err)
	}
	defer outputFile.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not receive chunk: %v", err)
		}

		if _, err := outputFile.Write(chunk.GetData()); err != nil {
			log.Fatalf("could not write to output file: %v", err)
		}
	}

	fmt.Println("File received successfully")
}
