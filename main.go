package main

import (
    al "github.com/azul3d/native-al"
    "log"
    "math/rand"
    "unsafe"
    "time"
)

const (
    freq = 44100 // 44.1 kHz
    duration = 3 // 3 seconds
)

func main() {
    device, err := al.OpenDevice("", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer device.Close()

    data := [freq*duration]int16{}
    for i := 0; i < freq*duration; i++ {
        data[i] = rnd(-32767, 32767);
    }
    var buffers uint32 = 0
    device.GenBuffers(1, &buffers)
    device.BufferData(buffers, al.FORMAT_MONO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*len(data)), freq)
    var sources uint32 = 0
    device.GenSources(1, &sources)
    device.Sourcei(sources, al.BUFFER, int32(buffers))
    device.SourcePlay(sources)

    time.Sleep(time.Second * 3)

    device.DeleteSources(1, &sources)
    device.DeleteBuffers(1, &buffers)
    log.Println("Done.")
}

func rnd(min, max int) int16 {
    return int16(min + (rand.Intn(max - min)))
}