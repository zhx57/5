package _plugin

import (
	"bytes"
	"testing"
)

// TestPbPageParseSynthetic 用合成 protobuf 验证 weltolkPbPageParse 的字段解析。
func TestPbPageParseSynthetic(t *testing.T) {
	// error: errorno=0
	errBuf := bytes.NewBuffer(nil)
	errBuf.Write(pbpEncodeInt32(1, 0))

	// page: total_page=2 (field 5)
	pageBuf := bytes.NewBuffer(nil)
	pageBuf.Write(pbpEncodeInt32(5, 2))

	// thread: replyNum=42 (field 4)
	threadBuf := bytes.NewBuffer(nil)
	threadBuf.Write(pbpEncodeInt32(4, 42))

	// content: text="hello" (field 2)
	contentBuf := bytes.NewBuffer(nil)
	contentBuf.Write(pbpEncodeString(2, "hello"))

	// author: name_show="Alice" (field 4)
	authorBuf := bytes.NewBuffer(nil)
	authorBuf.Write(pbpEncodeString(4, "Alice"))

	// post: id=100(1), floor=1(3), content(5), author_id=999(19), author(23)
	postBuf := bytes.NewBuffer(nil)
	postBuf.Write(pbpEncodeInt64(1, 100))
	postBuf.Write(pbpEncodeInt32(3, 1))
	postBuf.Write(pbpEncodeMessage(5, contentBuf.Bytes()))
	postBuf.Write(pbpEncodeInt64(19, 999))
	postBuf.Write(pbpEncodeMessage(23, authorBuf.Bytes()))

	// data: page(3), post_list(6), thread(8)
	dataBuf := bytes.NewBuffer(nil)
	dataBuf.Write(pbpEncodeMessage(3, pageBuf.Bytes()))
	dataBuf.Write(pbpEncodeMessage(6, postBuf.Bytes()))
	dataBuf.Write(pbpEncodeMessage(8, threadBuf.Bytes()))

	// response: error(1), data(2)
	respBuf := bytes.NewBuffer(nil)
	respBuf.Write(pbpEncodeMessage(1, errBuf.Bytes()))
	respBuf.Write(pbpEncodeMessage(2, dataBuf.Bytes()))

	result := weltolkPbPageParse(respBuf.Bytes())
	if result.errorno != 0 {
		t.Fatalf("errorno=%d, want 0", result.errorno)
	}
	if result.totalPage != 2 {
		t.Fatalf("totalPage=%d, want 2", result.totalPage)
	}
	if result.replyCount != 42 {
		t.Fatalf("replyCount=%d, want 42", result.replyCount)
	}
	if len(result.floors) != 1 {
		t.Fatalf("floors=%d, want 1", len(result.floors))
	}
	f := result.floors[0]
	if f.ID != 100 {
		t.Fatalf("ID=%d, want 100", f.ID)
	}
	if f.Floor != 1 {
		t.Fatalf("Floor=%d, want 1", f.Floor)
	}
	if f.Content != "hello" {
		t.Fatalf("Content=%q, want %q", f.Content, "hello")
	}
	if f.AuthorID != 999 {
		t.Fatalf("AuthorID=%d, want 999", f.AuthorID)
	}
	if f.Username != "Alice" {
		t.Fatalf("Username=%q, want %q", f.Username, "Alice")
	}
}

func TestPbPageRealRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real network test in short mode")
	}
	bduss := "10a0pQbkJ5a34wSUQ4QkpEZXFmcC1zZEVQSEdkS0xva1I5bW5janhBc3JlbHRxSVFBQUFBJCQAAAAAAAAAAAEAAAAG6LBSvaO2ybquva0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACvtM2or7TNqay"
	stoken := "d39f6e34073909bf939b4170d6d3a9674d68bb58b2406cdd728da93287b8785e"
	tid := int64(9534834391)
	replyCount, totalPage, ok := weltolkGetReplyCount(tid, bduss, stoken)
	if !ok {
		t.Fatalf("weltolkGetReplyCount failed")
	}
	if replyCount < 0 {
		t.Fatalf("replyCount=%d", replyCount)
	}
	floors := weltolkGetLastFloorContent(tid, bduss, stoken, 1, totalPage)
	if len(floors) == 0 {
		t.Fatalf("no floors")
	}
	t.Logf("replyCount=%d totalPage=%d floors=%d latest_floor=%d", replyCount, totalPage, len(floors), floors[0].Floor)
}
