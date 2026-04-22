// Guess a number, source unknown
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
int main()                           {
  int secret_num = 0, count = 0, num ;
  int stime                          ;
  long ltime                         ;
  ltime = time(NULL)                 ;
  stime = (unsigned)ltime / 2        ;
  srand(stime)                       ;
  secret_num = rand() % 100          ;

  while (1)                          {
    printf("\nGuess the number! ")   ;
    scanf("%d", &num)                ;
    if (secret_num == num)           {
      printf("You win!\n")           ;
      break                          ;}
    else if (secret_num < num)       {
      printf("Too big!")             ;}
    else if (secret_num > num)       {
      printf("Too small!")           ;}}
  return 0                           ;}
