Output "demo.gif"
Set Theme "Catppuccin Frappe"
Set Margin 24
Set MarginFill "#674EFF"
Set BorderRadius 20
Set Width 1200
Set Height 600
Set PlaybackSpeed 1.0

Type "qasm highlight examples/test.qasm"
Enter
Sleep 2s
Type "clear"
Sleep 1s
Enter

Type 'printf %b "OPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[2] q;\nbit[2] c;\nh q[0];\ncx q[0], q[1];\nmeasure q -> c;\n" | qasm fmt | qasm highlight'
Enter
Sleep 2s
Type "clear"
Enter
Sleep 1s

Type "qasm lint examples/test_bounds.qasm"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 1s
