package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "hoshiyomi/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn_creds grpc.DialOption = grpc.WithTransportCredentials(insecure.NewCredentials())

func printQuote(r *pb.SuiQuoteResp) {
	var quoteFooter string
	switch r.Source {
	case "":
		quoteFooter = "\033[34m- Hoshimachi Suisei\033[0m"
	default:
		quoteFooter = "\033[34m- Hoshimachi Suisei | " + r.Source + "\033[0m"
	}
	switch r.Type {
	case 0:
		fmt.Printf("\033[36m\033[22m \"%s\" \033[0m\n", r.Quote)
		fmt.Println(quoteFooter)
	case 1:
		fmt.Printf("\033[31m\033[22m \"%s\" \033[0m\n", r.Quote)
		fmt.Println(quoteFooter)
	}

}

func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Failed to load .env file - maybe it doesn't exist?")
	}

	sui_conn, conn_err := grpc.NewClient(os.Getenv("GRPC_DIAL_ADDRESS"), conn_creds)
	if conn_err != nil {
		log.Fatal("Failed to connect to gRPC server:", conn_err)
	}
	defer sui_conn.Close()

	client := pb.NewHoshimachiQuoteGenClient(sui_conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sui_resp, resp_err := client.GetRandomQuote(ctx, &pb.SuiQuoteReq{})
	if resp_err != nil {
		log.Fatal("Could not get response:", resp_err)
	}

	printQuote(sui_resp)
}
