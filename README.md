# enotes

View, write and edit encrypted notes

## Installation

```
go install github.com/zd4y/enotes@latest
```

## Usage

```
enotes
```

You will be prompted to enter a password that will be used for the encryption of the notes. After
that, you can create a new note pressing `enter` on `New note` in the list, this will open your
editor in a temporary file and after you quit the editor, the note will get encrypted and saved in
your current directory (the EDITOR environment variable is used here). Now you can see the note
contents pressing `enter` on the note title in the list and edit it pressing `e` in the note view.
