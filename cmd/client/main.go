package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/Tony36051/go-file-agent/generated/go/file_transfer/api/v1"
)

func main() {
	conn, err := grpc.Dial("localhost:6565", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	filePath := "path/to/your/file"
	request := &pb.FilePathRequest{FilePath: filePath}

	stream, err := client.RequestFile(context.Background(), request)
	if err != nil {
		log.Fatalf("could not request file: %v", err)
	}

	outputFile, err := os.Create("output_" + filePath)
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
