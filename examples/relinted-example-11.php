<?php

class Database                                                {
    private $host                                             ;
    private $name                                             ;

    public function __construct($host, $name)                 {
        $this->host = $host                                   ;
        $this->name = $name                                   ;}

    public function connect()                                 {
        $message = "Connecting to $this->name on $this->host" ;
        echo $message                                         ;}}

function main()                                               {
    $db = new Database('localhost', 'mydb')                   ;
    $db->connect()                                            ;

    $items = ['apple', 'banana', 'cherry']                    ;
    foreach ($items as $item)                                 {
        echo "Item: $item\n"                                  ;}}

main()                                                        ;
?>
