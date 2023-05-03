# Maelstrom sounds

The `snd_*.mp3` files in this directory originate from Maelstrom by
[Ambrosia Software](https://en.wikipedia.org/wiki/Ambrosia_Software).

From Wikipedia:

> Ambrosia Software was a predominantly Macintosh software company
> founded in 1993 and located in Rochester, New York, U.S.
> Ambrosia Software was best known for its Macintosh remakes of older
> arcade games, which began with a 1992 version of Atari, Inc.'s Asteroids from 1979.
> ...
> The first game distributed under the Ambrosia Software name was Maelstrom, a 1992
> remake of the 1979 Asteroids arcade video game

A Open-Source fork of Maelstrom luckily exists on [Github](https://github.com/richardjs/Maelstrom).

From its [README](https://github.com/richardjs/Maelstrom/blob/master/README.md):

> Maelstrom was originally a Mac OS game, first released as shareware
> in 1992 by Andrew Welch of Ambrosia Software. Ambrosia gave the source
> to Sam Lantinga, who released Maelstrom 3.0, a SDL port under the GPL, in 1995.
> The game assets were released under the CC Attribution license in 2010.

To have some audio feedback within k12-booter, I first thought of
Maelstrom sounds, as it was actually the first MacOS game I ever
played thanks to Lothar Zeidler in New Jersey. To conclude, the
OSS Maelstrom repository is awesome to study ... and/or despair ;).
Files like [Mac_Wave.cpp](https://github.com/richardjs/Maelstrom/blob/master/maclib/Mac_Wave.cpp#L554):

```cpp
	/* Figure out how big the RIFF chunk will be */
	wavelen = sizeof(Uint32)+sizeof(format)+2*sizeof(Uint32)+sound_datalen;

	/* Save the WAVE */
```

That's enough for me. So I just ran the `snd2wav` tool to extract the Maelstrom
samples from the [Maelstrom_Sounds](https://github.com/richardjs/Maelstrom/blob/master/Maelstrom_Sounds)
resource fork like this on my mac:

```bash
# prerequisites via https://brew.sh
brew install sdl sdl2 sdl_net sdl12-compat

git clone https://github.com/richardjs/Maelstrom.git
cd Maelstrom
./configure
make

cd maclib
./snd2wav -rate 44100 ../Maelstrom_Sounds
for s in *.wav; do ffmpeg -i $s $s.mp3; done
```

See also [Maelstrom gameplay on YouTube](https://www.youtube.com/watch?v=bWGI-MHSpH8).

# PEACE

Free Julian Assange.

