<!doctype html>
<html>
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <link href="/tailwind.css" rel="stylesheet">
      <title></title>
    </head>
    <body>
        <div class="text-xl font-medium text-black">
            <form action="/" method="post">
                Miner ID:<input type="text" name="miner-id">
                Start Date:<input type="date" name="start-date" value={{.Start}}>
                End Date:<input type="date" name="end-date" value={{.End}}>
                <input type="submit" value="Go">
            </form>
        </div>
    </body>
</html>
