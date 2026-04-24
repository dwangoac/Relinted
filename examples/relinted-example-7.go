package main //Source: https://freshman.tech/golang-guess
import (
        "bufio"
        "fmt"
        "math/rand"
        "os"
        "strconv"
        "strings"
        "time"
)

func main()                                               {
        min, max := 1, 100
        rand.Seed(time.Now().UnixNano())
        secretNumber := rand.Intn(max-min) + min
        fmt.Println("Guess the number!")

        for                                               {
               reader := bufio.NewReader(os.Stdin)
               input, err := reader.ReadString('\n')
               if err != nil                              {
                      fmt.Println("Failed to read", err)
                      continue                            }

               input = strings.TrimSuffix(input, "\n")
               guess, err := strconv.Atoi(input)
               fmt.Println("You Guessed:", guess)

               if guess > secretNumber                    {
                      fmt.Println("Too big!")             }
               else if guess < secretNumber               {
                      fmt.Println("Too small!")           }
               else                                       {
                      fmt.Println("You win!")
                      break                               }}}
