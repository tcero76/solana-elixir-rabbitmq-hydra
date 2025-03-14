package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

type TransferEvent struct {
	Sender   solana.PublicKey `json:"sender"`
	Receiver solana.PublicKey `json:"receiver"`
	Amount   uint64           `json:"amount"`
}

func main() {

	ctx := context.Background()
	// URL del WebSocket de Solana Devnet (cambiar si es otro cluster)
	rpcWsUrl := "wss://api.devnet.solana.com"

	// Conectarse al WebSocket
	client, err := ws.Connect(context.Background(), rpcWsUrl)
	if err != nil {
		log.Fatalf("Error al conectar al WebSocket: %v", err)
	}
	defer client.Close()

	// Dirección del programa (reemplázala con la de tu contrato)
	programID := os.Getenv("PROGRAMID")
	pubKey, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		log.Fatalf("Error al convertir programID a PublicKey: %v", err)
	}
	// Suscribirse a los logs de transacciones donde participa el programa
	sub, err := client.LogsSubscribeMentions(
		pubKey,
		rpc.CommitmentProcessed,
	)
	if err != nil {
		log.Fatalf("Error al suscribirse a los logs: %v", err)
	}
	defer sub.Unsubscribe()

	// Conectar a RabbitMQ
	rabbitMQUrl := os.Getenv("BROKER")
	conn, err := amqp091.Dial(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al crear canal de RabbitMQ: %v", err)
	}
	defer ch.Close()
	fmt.Println("📡 Escuchando eventos en la red...")
	// Declarar una cola para recibir los mensajes
	queue, err := ch.QueueDeclare(
		"solana_event_queue", // nombre de la cola
		true,                 // durable
		false,                // auto-delete
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Error al declarar cola en RabbitMQ: %v", err)
	}
	// Función para enviar el mensaje a RabbitMQ
	sendToRabbitMQ := func(message string) {
		err := ch.Publish(
			"",         // exchange
			queue.Name, // routing key (cola)
			false,      // mandatory
			false,      // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			log.Fatalf("Error al enviar mensaje a RabbitMQ: %v", err)
		}
		fmt.Println("Mensaje enviado a RabbitMQ:", message)
	}
	// Leer eventos en un bucle
	for {
		event, err := sub.Recv(ctx)
		if err != nil {
			log.Fatalf("Error recibiendo logs: %v", err)
		}

		// Buscar logs que comiencen con "Program log: "
		for _, logEntry := range event.Value.Logs {
			if len(logEntry) > 12 && logEntry[:12] == "Program log: " {
				eventData := logEntry[12:]
				// Intentar decodificarlo como TransferEvent
				var event TransferEvent
				err := json.Unmarshal([]byte(eventData), &event)
				if err == nil {
					fmt.Printf("🟢 Evento recibido: %+v\n", event)
				}
			}
		}

		sendToRabbitMQ(strings.Join(event.Value.Logs, ","))
		fmt.Println("🔔 Nuevo evento recibido:")
		fmt.Println(event.Value.Logs) // Aquí estarán los logs del evento TransferEvent
	}
}
