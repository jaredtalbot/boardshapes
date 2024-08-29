echo Downloading Godot...
curl -o godot.zip -L 'https://github.com/godotengine/godot/releases/download/4.3-stable/Godot_v4.3-stable_linux.x86_64.zip'
echo Installing unzip...
sudo apt install unzip
echo Unzipping Godot...
unzip -p godot.zip Godot_v*-stable_linux.x86_64 > godot.x86_64
chmod +xr godot.x86_64
echo Checking files...
find . -name '*.gd' -print0 | xargs -P 8 -n 1 -0  ./godot.x86_64 --headless --check-only -q -s
if [ $? -eq 0 ]
then
    echo All files checked successfully.
    exit 0
else 
    echo There is an issue with one of the files.
    exit 1
fi
