package handlers

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/Tony36051/go-file-agent/generated/go/file_transfer/api/v1"
)

// DownloadFile handles file download requests by calling the gRPC server
func DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial("127.0.0.1:6565", grpc.WithInsecure())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	request := &pb.FilePathRequest{FilePath: filename}

	stream, err := client.RequestFile(context.Background(), request)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Set response headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// Stream the file content
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			if _, err := c.Writer.Write(chunk.GetData()); err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
	}()
	wg.Wait()
}
