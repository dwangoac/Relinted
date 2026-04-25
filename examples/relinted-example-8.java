package com.example                                ;

import java.util.Scanner                           ;

/**
 * Simple calculator demo.
 * Demonstrates Java formatting.
 */
public class Calculator                            {
    private int result = 0                         ;
    private String label = "Calculator"            ;

    public void add(int n)                         {
        result += n                                ;
        /* Update result */                        }

    public void subtract(int n)                    {
        result -= n                                ;}

    public int getResult()                         {
        return result                              ;}

    public char getSeparator()                     {
        char sep = '|'                             ;
        return sep                                 ;}

    public static void main(String[] args)         {
        Calculator calc = new Calculator()         ;
        calc.add(10)                               ;
        calc.add(5)                                ;

        String help = """
            Usage: add <n>, subtract <n>, result
            """                                    ;

        if (calc.getResult() > 0)                  {
            System.out.println("Positive!")        ;}
        else                                       {
            System.out.println("Zero or negative") ;}}}
