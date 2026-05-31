# Small tui calendar for Omarchy/Waybar

[Omarchy](https://omarchy.org/) clock widget in waybar by default switches 
time format. I want it to show a small lightweight calendar instead, so I made 
this.



https://github.com/user-attachments/assets/4cd55de1-3694-413b-8cd8-15315deba2d0



I don't want it to be a full-fledged web app, so it's 
a [bubbletea](https://github.com/charmbracelet/bubbletea) TUI.

It's as simple as it gets, as most likely this all will be obsolete in 
Omarchy 4, with the switch to Quickshell. 

No events, ICS, etc. Just displays a month grid and the current date.
Arrow or hjkl navigation for moving the selected date around; 
shift+{left/right/h/l} to switch months around.

q,esc,ctrl+c to exit.

## Installation:

You need go toolchain installed. (Install -> Development -> Go)

Install the binary using go itself:

```sh
go install github.com/religiosa1/omacal@latest
```

Modify your waybar config (`~/.config/waybar/config.jsonc`), to use this binary
on click, e.g.:

```jsonc
...
"clock": {
  "format": "{:L%A %H:%M}", // or whatever you like, I use "W{:L%W %d %a %H:%M}"
  // comment out alternative format
  // "format-alt": "{:L%d %B W%V %Y}", 
  "format": "W{:L%W %d %a %H:%M}",
  "tooltip": false,
  // And add this line:
  "on-click": "omarchy-launch-or-focus-tui omacal",
  "on-click-right": "omarchy-launch-floating-terminal-with-presentation omarchy-tz-select"
},
...
```

Restart your waybar:

```sh
pkill waybar && hyprctl dispatch exec waybar
```

And add some floating window styles in hypr (`~/.config/hypr/hyprland.conf`):

```
windowrule {
  name = omacal
  match:class = (org.omarchy.omacal)
  float = on
  move = (monitor_w*0.50-115) (32)
  size = 230 174
  border_size = 1
}
```

If you're playing around with styles, remember a month can contain up to 6 rows 
of weeks in a calendar (e.g. March 2026).

## License

Omacal is MIT licensed.
