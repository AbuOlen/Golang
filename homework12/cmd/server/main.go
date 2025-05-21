package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	cmds "hw12/internal/commands"
	store "hw12/internal/documentstore"
)

const primaryKey = "key"
const collectionKey = "key"
var s = store.NewStore()

func execPut(raw string, col *store.Collection) (string, error) {
	p := &cmds.PutCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}
	d1 := store.Document{Fields: make(map[string]store.DocumentField)}
	d1.Fields[primaryKey] = store.DocumentField{Type: store.DocumentFieldTypeString, Value: p.Key}
	d1.Fields["val"] = store.DocumentField{Type: store.DocumentFieldTypeString, Value: p.Value}
	col.Put(d1)

	resp := &cmds.PutCommandResponsePayload{}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	fmt.Println(rawResp)
	return string(rawResp), nil
}

func execGet(raw string, col *store.Collection) (string, error) {
	p := &cmds.GetCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	doc, ok := col.Get(p.Key)
	var value string
	if ok {
		value = doc.Fields["val"].Value.(string)
	}

	resp := &cmds.GetCommandResponsePayload{
		Value: value,
		Ok:    ok,
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execDelete(raw string, col *store.Collection) (string, error) {
	p := &cmds.DeleteCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	ok := col.Delete(p.Key)
	resp := &cmds.DeleteCommandResponsePayload{
		Ok: ok,
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execList(col *store.Collection) (string, error) {
	documents := col.List()
	values := make([]string, len(documents))

	for i, doc := range documents {
		values[i] = doc.Fields["val"].Value.(string)
	}

	resp := &cmds.ListCommandResponsePayload{
		Value: values,
		Ok:    true,
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func handleConnection(conn net.Conn, key string) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	w := bufio.NewWriter(conn)

	for scanner.Scan() {
		msg := scanner.Text()

		elems := strings.Split(msg, " ")
		length := len(elems)
		if length != 1 && length != 2 {
			w.WriteString("invalid command\n")
			w.Flush()
			continue
		}

		col, found := s.GetCollection(key)
		if !found {
			panic("Collection not found")
		}
		var resp string
		var err error

		switch elems[0] {
		case cmds.PutCommandName:
			resp, err = execPut(elems[1], col)
		case cmds.GetCommandName:
			resp, err = execGet(elems[1], col)
		case cmds.DeleteCommandName:
			resp, err = execDelete(elems[1], col)
		case cmds.ListCommandName:
			resp, err = execList(col)
		default:
			w.WriteString("invalid command\n")
		}

		if err != nil {
			w.WriteString(fmt.Sprintf("error: %s\n", err))
		}

		w.WriteString(fmt.Sprintf("response: %s\n", resp))

		w.Flush()
	}

	fmt.Println("connection closed")
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(fmt.Errorf("error listening: %w", err))
	}

	cfg := store.CollectionConfig{PrimaryKey: primaryKey}
	ok, _ := s.CreateCollection(collectionKey, &cfg)
	if !ok {
		fmt.Println("Collection creation failed")
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(fmt.Errorf("error accepting connection: %w", err))
			return
		}

		fmt.Println("connection accepted")

		go handleConnection(conn, collectionKey)
	}
}