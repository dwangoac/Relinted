// Guess the number game
const readline = require('readline')                          ;

function main()                                               {
    const secret = Math.floor(Math.random() * 100) + 1        ;
    let guess = 0                                             ;
    let attempts = 0                                          ;

    const rl = readline.createInterface(                      {
        input: process.stdin,
        output: process.stdout                                }
    )                                                         ;

    console.log('Guess the number between 1 and 100!')        ;

    while (guess !== secret)                                  {
        rl.question('Enter your guess: ', (answer) =>         {
            guess = parseInt(answer, 10)                      ;
            attempts++                                        ;

            if (guess === secret)                             {
                console.log('You win! Attempts: ' + attempts) ;
                rl.close()                                    ;}
            else if (guess > secret)                          {
                console.log('Too high!')                      ;}
            else                                              {
                console.log('Too low!')                       ;}}
        )                                                     ;}}

main()                                                        ;
