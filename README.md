# filacct

Track mining accounts

## Track a single Miner Actor

```
$ filacct download <miner_actor id>
```

## Track all the Miner Actors of a Owner Address

```
$ filacct download --owner <owner account_actor id>
```

Note: The `download` command should be run by a cronjob/systemd-timer to remanin in sync.

## Serve reporting frontend

Run the following from the same directory as the `download` commands were run:
```
# filacct serve --port=80
```

Note: Needs to be run as root if you don't want to setup a reverse proxy application.

## Credits

Initial code written by [Factor*8 Solutions](https://github.com/Factor8Solutions/li_fil-qnd-burn-sheet).
