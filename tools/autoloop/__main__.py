import argparse
import librosa.core
import mutagen
from pymusiclooper.core import MusicLooper


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("file")
    args = parser.parse_args()

    looper = MusicLooper(args.file)
    loop_pair = looper.find_loop_pairs()[0]
    loop_start = looper.frames_to_samples(loop_pair["loop_start"])
    loop_end = looper.frames_to_samples(loop_pair["loop_end"])
    loop_length = loop_end - loop_start

    print(
        f'found loop with score {loop_pair["score"]*100:4.2f}% probability: [{loop_start}, {loop_end})'
    )
    print(
        f'in milliseconds: [{librosa.core.frames_to_time(loop_pair["loop_start"] * 1000, sr=looper.rate)}, {librosa.core.frames_to_time(loop_pair["loop_end"] * 1000, sr=looper.rate)})'
    )
    mf = mutagen.File(args.file)
    mf["LOOPSTART"] = str(loop_start)
    mf["LOOPLENGTH"] = str(loop_length)
    mf.save()


if __name__ == "__main__":
    main()
