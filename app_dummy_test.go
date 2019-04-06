package lazyledger

import (
    "testing"
)


func TestAppDummySimpleBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewDummyApp(ms)
    b.RegisterApplication(&app)

    puts := make(map[string]string)
    puts["foo"] = "bar"
    puts["goo"] = "tar"

    sb.AddMessage(app.(*DummyApp).GenerateTransaction(puts))
    b.ProcessBlock(sb)

    if app.(*DummyApp).Get("foo") != "bar" || app.(*DummyApp).Get("goo") != "tar" {
        t.Error("dummy app state update failed")
    }
}
