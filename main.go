package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var (
    resolution string // Variable to hold the resolution flag

    downloadCmd = &cobra.Command{
        Use:   "download [url]",
        Short: "Download media from supported platforms",
        Long:  "Download videos and photos from YouTube, Facebook, and Instagram.",
        Args:  cobra.ExactArgs(1),
        Run:   downloadMedia,
    }

    audioCmd = &cobra.Command{
        Use:   "audio [url]",
        Short: "Convert YouTube video to audio",
        Long:  "Download and convert a YouTube video to audio format.",
        Args:  cobra.ExactArgs(1),
        Run:   convertToAudio,
    }

    listFormatsCmd = &cobra.Command{
        Use:   "list-formats [url]",
        Short: "List available formats for a video",
        Long:  "List all available formats for the specified video URL.",
        Args:  cobra.ExactArgs(1),
        Run:   listFormats,
    }
)

func main() {
    var rootCmd = &cobra.Command{Use: "media-cli"}

    // Add resolution flag to the download command
    downloadCmd.Flags().StringVarP(&resolution, "resolution", "r", "best", "Specify the maximum resolution (e.g., 720, 1080)")

    // Add commands to the root
    rootCmd.AddCommand(downloadCmd, audioCmd, listFormatsCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func downloadMedia(cmd *cobra.Command, args []string) {
    url := args[0]
    downloadDir := "downloads"

    // Create the download directory if it doesn't exist
    os.MkdirAll(downloadDir, os.ModePerm)

    // Convert resolution string to integer
    resInt, err := strconv.Atoi(resolution)
    if err != nil {
        fmt.Println("Invalid resolution format. Please enter a valid number (e.g., 720, 1080).")
        return
    }

    // Use yt-dlp to download the media with the specified resolution and attach audio
    cmdStr := []string{"-f", fmt.Sprintf("bestvideo[height<=%d]+bestaudio/best", resInt), "-o", filepath.Join(downloadDir, "%(title)s.%(ext)s"), url}
    
    err = executeCommand("yt-dlp", cmdStr)
    
    if err != nil {
        fmt.Println("Error downloading media:", err)
        fmt.Println("You can use 'media-cli list-formats [url]' to see available formats.")
        return
    }

    fmt.Println("Download completed successfully.")
}

func convertToAudio(cmd *cobra.Command, args []string) {
    url := args[0]
    audioDir := "audio"

    // Create the audio directory if it doesn't exist
    os.MkdirAll(audioDir, os.ModePerm)

    // Use yt-dlp to download the audio
    cmdStr := []string{"--extract-audio", "--audio-format", "mp3", "-o", filepath.Join(audioDir, "%(title)s.%(ext)s"), url}
    err := executeCommand("yt-dlp", cmdStr)

    if err != nil {
        fmt.Println("Error downloading audio:", err)
        return
    }

    fmt.Println("Audio downloaded successfully.")
}

func listFormats(cmd *cobra.Command, args []string) {
    url := args[0]
    
    // List available formats using yt-dlp
    cmdStr := []string{"-F", url}
    err := executeCommand("yt-dlp", cmdStr)
    if err != nil {
        fmt.Println("Error listing formats:", err)
    }
}

func executeCommand(command string, args []string) error {
    // Create the command
    cmd := exec.Command(command, args...)

    // Set output to standard
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to execute command: %w", err)
    }
    return nil
}
