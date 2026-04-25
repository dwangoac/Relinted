using System;
using System.Collections.Generic;

namespace Demo
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var names = new List<string> { "Alice", "Bob", "Charlie" };

            foreach (var name in names)
            {
                Console.WriteLine("Hello, " + name);
            }

            int sum = 0;
            for (int i = 0; i < 5; i++)
            {
                sum += i;
            }

            Console.WriteLine("Sum: " + sum);
        }
    }
}
