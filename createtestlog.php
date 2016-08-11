<?php
//
// Generate test log via `php createtestlog.php > irctest.log`
//

date_default_timezone_set("GMT");
$begin = new DateTime();
$begin = $begin->modify('-720 day');
$end = new DateTime();

for($i = $begin; $i <= $end; $i->modify('+1 day')){

    $numberOfLines = rand(0, 250);
    $linesGenerated = 0;
    $time = new DateTime("2000-01-01 00:00:00");

    $people = [
        'joeblogs',
        'billythekid',
        'jasper',
        'lelu',
        'tron',
        'keeper',
        'jullia',
        'teddy',
        'jon',
        'zzt',
        'nation',
        'ms',
        'debs',
        'the-don'
    ];

    while ($linesGenerated < $numberOfLines) {

        $username = $people[array_rand($people, 1)];
        if (rand(0,1) == 1){
            $message = "<$username> This is a message";
        }else{
            $message = "* $username is preforming an action";
        }

        $time->modify('+'. rand(15,600) .' second');
        echo "[" . $i->format("Y-m-d ") . $time->format("H:i:s O") . "] " . $message . PHP_EOL;
        $linesGenerated++;
    }
}