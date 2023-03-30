# TODO 

Once HTML template parsing works, provide cui code generation templates for ...

1. 1st form: select country, confirm default buckets / api endpoint
1. 2nd form: school base data or skip to next
1. 3rd form: request school ID (grpc/api request) / provide it via form
1. 4th form: add hardware inventory items (xml files? api request)
1. 5th form: table view of hardware inventory
1. 6th form: configure OSS sources, releases
1. 7th form: Link farm, nested: Deployment, Monitoring, Services ...

Once all form stuff is implemented as cui template,
remove any hand-written cui form code and provide
HTML templates to reconstruct it.

Goals:

1. Maintain HTML FORMs as smallest denominator in git as source for UI code generation
1. improvements in formgenerator should allow regeneration of cui/ui code (don't touch generated code)
1. reliance on few basic elements ensures ability to provide *consistent* UI / UX

Assumption: I look at my screen to deal with ice cold data. Roses or so distract me in many ways here. AS400 FTW.

# Random GUI detail thoughts

- command-"line" (rows/scrollback!) could provide access to built-in functions; like system() m(
- switch between command-line and stack view (stack being the better clipboard)
- allow command-line stack access
- show nerdy status info like terminal emulator: CAPS_LOCK RXTX CLOCK LOCALE
- introduce color scheme before it hurts to do so
- given standard UIs SHOULD make this even easier, look into braille support; text-to-speech?
