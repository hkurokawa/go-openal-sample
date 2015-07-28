package main

import (
    al "azul3d.org/native/al.v1"
    audio "azul3d.org/audio.v1"
    _ "azul3d.org/audio/wav.v1"
    "log"
    "math/rand"
    "unsafe"
    "time"
    "os"
)

const (
    freq = 44100 // 44.1 kHz
    duration = 3 // 3 seconds
)

func main() {
    var data []int16
    switch len(os.Args) {
        case 1:
        log.Println("Loading white noise.")
        data = genWhiteNoise()
        break

        case 2:
        filename := os.Args[1]
        log.Printf("Loading %s\n", filename)
        data = readFile(filename)
        break

        default:
        log.Fatalf("Unexpected number of arguments: %s. Must be 1 or 2.\nUsage: go run main.go [file]\n", len(os.Args))
    }


    device, err := al.OpenDevice("", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer device.Close()

    var buffers uint32 = 0
    device.GenBuffers(1, &buffers)
    device.BufferData(buffers, al.FORMAT_MONO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*len(data)), freq)
    var sources uint32 = 0
    device.GenSources(1, &sources)
    device.Sourcei(sources, al.BUFFER, int32(buffers))
    device.SourcePlay(sources)

    time.Sleep(time.Second * duration)

    device.DeleteSources(1, &sources)
    device.DeleteBuffers(1, &buffers)
    log.Println("Done.")
}

func readFile(filename string) []int16 {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }

    // Create a decoder for the audio source
    decoder, format, err := audio.NewDecoder(file)
    if err != nil {
        log.Fatal(err)
    }
    config := decoder.Config()
    log.Printf("Decoding a %s file.\n", format)
    log.Println(config)

    // Create a buffer that can hold 3 second of audio samples
    bufSize := duration * config.SampleRate * config.Channels
    buf := make(audio.F64Samples, bufSize)

    // Fill the buffer with as many audio samples as we can
    read, err := decoder.Read(buf)
    if err != nil {
        log.Fatal(err)
    }
    readBuf := buf.Slice(0, read)
    data := make([]int16, readBuf.Len())
    for i := 0; i < readBuf.Len(); i++ {
        data[i] = int16(audio.F64ToPCM16(readBuf.At(i)))
    }
    return data
}

func genWhiteNoise() []int16 {
    data := make([]int16, freq*duration)
    for i := 0; i < freq*duration; i++ {
        data[i] = rnd(-32767, 32767);
    }
    return data
}

func rnd(min, max int) int16 {
    return int16(min + (rand.Intn(max - min)))
}