# Save and plot your typeracer stats

## Configure

Add your username and password in the configuration file.

```sh
cp config.yml.example config.yml
vim config.yml
```

## Build

```sh
go build -v -o typeracer .
```

## Save your stats in a json file

The stats will be saved in `stats.json`.

```sh
./typeracer save
```

## Generate the plots

The plot will be saved in `plot.png`

```sh
./typeracer plot
```
