package main

// import "mygo/hello"
import (
    "fmt"
    "reflect"
)

func main() {
    text1 := "Go語言"
    for i := 0; i < len(text1); i++ {
        fmt.Printf("%x ", text1[i])
    }
    fmt.Println()

    for index, runeValue := range text1 {
        fmt.Printf("%#U 位元起始位置 %d\n", runeValue, index)
    }
    fmt.Println()

    bs := []byte(text1)
    bs[0] = 103
    text2 := string(bs)
    fmt.Println(text1) // Go語言
    fmt.Println(text2) // go語言


    s1 := []int{1, 2, 3, 4, 5}
    a1 := [...]int{1, 2, 3, 4, 5}
    fmt.Println(reflect.TypeOf(s1)) // []int
    fmt.Println(reflect.TypeOf(a1)) // [5]int


    arr := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    slice := arr[3:4]
    fmt.Println(reflect.TypeOf(arr))   // [10]int
    fmt.Println(reflect.TypeOf(slice)) // []int
    fmt.Println(len(slice))            // 3
    fmt.Println(cap(slice))            // 9

    fmt.Println(slice)   // [2 3 4]
    fmt.Println(arr)     // [1 2 3 4 5 6 7 8 9 10]

    slice[0] = 20
    fmt.Println(slice)   // [20 3 4]
    fmt.Println(arr)     // [1 20 3 4 5 6 7 8 9 10]
}

