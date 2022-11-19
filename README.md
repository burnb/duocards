## DuoCards CLI

Additional functionality for the [DuoCards.com](https://duocards.com)

* Cards import from file

### Environment vars

```
API_TOKEN={DuoCards API Token}
DECK_ID={DuoCards Deck ID}
```

Also it's just possible to create the .env file near your binary(look at .env.example file)

### Usage

```
Usage: duocards [OPTIONS] COMMAND

Options:
  -h, --help            Show this help message and exit
  -i INPUT_FILE         Input file path (default "text.txt")
  -l LANGUAGE           Your native language (default "en")
  -s FILE_SEPARATOR     Columns separator in the input file (default "/")
  
Commands:
    import              Import from file
```