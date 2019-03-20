package audio

import (
	"encoding/binary"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
)

func PlaybackAudio(audioFileName string, guildId string, channelId string, session *discordgo.Session) (err error) {
	audioData, err := loadAudioData(audioFileName)
	if err != nil {
		return err
	}

	voiceChannel, err := session.ChannelVoiceJoin(guildId, channelId, false, false)
	if err != nil {
		return err
	}

	defer func() { voiceChannel.Disconnect() }()

	voiceChannel.Speaking(true)
	for _, buff := range audioData {
		voiceChannel.OpusSend <- buff
	}
	voiceChannel.Speaking(false)

	return nil
}

func Exists(audioFileName string) bool {
	_, err := os.Stat("assets/audio/" + audioFileName + ".dca")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func loadAudioData(audioFileName string) (audioData [][]byte, err error) {
	audioFile, err := os.Open("assets/audio/" + audioFileName + ".dca")
	if err != nil {
		return nil, err
	}

	var opusFrameLength int16
	audioData = make([][]byte, 0)

	for {
		err = binary.Read(audioFile, binary.LittleEndian, &opusFrameLength)

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				err = audioFile.Close()
				if err != nil {
					return nil, err
				}
				return audioData, nil
			}
			return nil, err
		}

		audioFrame := make([]byte, opusFrameLength)
		err = binary.Read(audioFile, binary.LittleEndian, &audioFrame)
		if err != nil {
			return nil, err
		}
		audioData = append(audioData, audioFrame)
	}
}
