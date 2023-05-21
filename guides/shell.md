# Shell

This guide aims in setting up YOUR shell. Be it bash, zsh, fish, ...
We always want some (same) toolset on our current box - and this
guide shall ideally document your clicks so you can simply re-apply
on your next workstation.

For developers and others seeking fancy git prompts, one common
requirements are #Nerdfonts today. Select a single one you like
or install the set for others to choose one. #CLICK#

Starship makes the terminal prompt fly - https://starship.rs/ #CLICK#
`starship config` this:
```
right_format = "$time"

[time]
disabled = false
style = "bold bright-black"
format = "[$time]($style)"

[line_break]
disabled = true
```
https://twitter.com/StarshipPrompt/status/1440380147807318020

# shell tools

grep is great but ripgrep is greater.
https://github.com/BurntSushi/ripgrep

