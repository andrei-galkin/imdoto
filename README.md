# Intro
 **imdoto** is a golang program to search keywords or phrases on Google Images Search and automatically download images to your computer. 
 
# Disclaimer
This console utility lets you download dozens of images from Google search engine. No copyright infringement is intended. We do not own the pictures which might be download by the tool. The tool is not meant to violate any licenses nor is any disrespect intended. Please do not download or use any image that violates its copyright terms. Google Images is a search engine that merely indexes images and allows you to find them. The tool does NOT produce its own images and, as such, it doesn't own copyright on any of them. The original creators of the images own the copyrights. Images published in the United States are automatically copyrighted by their owners, even if they do not explicitly carry a copyright warning. You may not reproduce copyright images without their owner's permission, except in "fair use" cases, or you could risk running into lawyer's warnings, cease-and-desist letters, and copyright suits. Please be very careful before its usage!

# Build
Go to **src** folder and run the command 
``` go build imdoto.go ```

# Arguments
Argument | Description
------------ | -------------
key   | Denotes the term or phrases you want to search for. For more than one keywords, wrap it in single quotes.
folder | The name of the folder where the downloaded images will be stored
type   | File types such as jpeg, png, bmp, gif or * is any file type
limit  | How many images you are going to download from the seach result

# Example
```imtodo -key "apple seed" -folder img -type * -limit 10```

# Contribute
Anyone is welcomed to contribute to this script. If you would like to make a change, open a pull request. For issues and discussion visit the Issue Tracker.
