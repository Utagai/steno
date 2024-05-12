# NOTE: This tape file is not "complete". VHS doesn't currently have a way of
# displaying keypresses on the screen, so what I've done is created the
# recording video (.mp4), and then edited the video in Da Vinci to add keypress
# info, and created a new .mp4 file. That new file is then converted to a GIF
# via ffmpeg & gifski and put into the repo and linked to from the README.
Output rsrc/demo.mp4

Require jam
Require gum

Set Shell "bash"
Set FontSize 24
Set Width 1200
Set Height 600

Type "jam"
Sleep 500ms
Enter
Sleep 2s
Type "b"
Sleep 1s
Type "d"
Sleep 1s
Type "?"
Sleep 3s
Backspace
Sleep 1s
Type "?"
Sleep 1s
Type "."
Sleep 5s 
Type "jam"
Sleep 300ms 
Enter
Sleep 1500ms 
Type "r"
Sleep 2s