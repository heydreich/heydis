package parser

import (
	"bytes"
	"fmt"
	resp "godis/interface/redis"
	"godis/lib/utils"
	reply "godis/redis/protocol"
	"io"
	"testing"
)

func TestParseStream(t *testing.T) {
	replies := []resp.Reply{
		reply.MakeIntReply(1),
		reply.MakeStatusReply("OK"),
		reply.MakeErrReply("ERR unknown"),
		reply.MakeBulkReply([]byte("a\r\nb")), // test binary safe
		reply.MakeNullBulkReply(),
		reply.MakeMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		reply.MakeEmptyMultiBulkReply(),
	}
	reqs := bytes.Buffer{}
	for _, re := range replies {
		reqs.Write(re.ToBytes())
		fmt.Println(re.ToBytes())
	}
	expected := make([]resp.Reply, len(replies))
	copy(expected, replies)

	reqs.Write([]byte("set a a")) // test text reply
	expected = append(expected, reply.MakeMultiBulkReply([][]byte{
		[]byte("set"), []byte("a"), []byte("a"),
	}))

	ch := ParseStream(bytes.NewReader(reqs.Bytes()))
	i := 0
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF {
				return
			}
			t.Error(payload.Err)

			return
		}
		if payload.Data == nil {
			t.Error("empty data")
			return
		}
		exp := expected[i]
		i++
		if !utils.BytesEquals(exp.ToBytes(), payload.Data.ToBytes()) {
			t.Error("parse failed: " + string(exp.ToBytes()))
		}
	}
}

func TestParseOne(t *testing.T) {
	replies := []resp.Reply{
		reply.MakeIntReply(1),
		reply.MakeStatusReply("OK"),
		reply.MakeErrReply("ERR unknown"),
		reply.MakeBulkReply([]byte("a\r\nb")), // test binary safe
		reply.MakeNullBulkReply(),
		reply.MakeMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		reply.MakeEmptyMultiBulkReply(),
	}
	for _, re := range replies {
		result, err := ParseOne(re.ToBytes())
		if err != nil {
			t.Error(err)
			continue
		}
		if !utils.BytesEquals(result.ToBytes(), re.ToBytes()) {
			t.Error("parse failed: " + string(re.ToBytes()))
		}
	}
}
