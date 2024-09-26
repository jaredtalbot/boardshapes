echo Downloading Godot...
curl -o godot.zip -L 'https://github.com/godotengine/godot/releases/download/4.3-stable/Godot_v4.3-stable_linux.x86_64.zip'
echo Installing unzip...
sudo apt install unzip
echo Unzipping Godot...
unzip -p godot.zip Godot_v*-stable_linux.x86_64 > godot.x86_64
chmod +xr godot.x86_64
echo Running project...
error=$(./godot.x86_64 --headless --import --path ./Game &> >(grep 'SCRIPT ERROR' - | wc -c))
if [ $error -gt 0 ]
then
    echo There is a script error.
    exit 1
else 
    echo Godot check successful.
    exit 0
fi
