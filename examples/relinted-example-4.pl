#!/usr/bin/perl
use 5.010                      ;

my $found  = 0                 ;
my $hidden = 1 + int rand 100  ;

while ( $found == 0 )          {
    print "Guess the number! " ;
    my $guess = <STDIN>        ;
    chomp $guess               ;

    if ( $guess < $hidden )    {
        say "Too small!"       ;}

    if ( $guess > $hidden )    {
        say "Too big!"         ;}

    if ( $guess == $hidden )   {
        $found = 1             ;
        say "You win!"         ;}}
