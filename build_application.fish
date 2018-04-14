#!/usr/bin/env fish
echo 'Building UI (Angular Web App)'
cd ui

echo Deleting 'dist' folder
rm -r dist

echo Compiling UI
ng build --prod --environment=prod

cd ../server/ui-dist
# Remove Old Binaries
echo "Removing old binaries"
for file in *
    set VALID 0
    for valid_file in (git ls-files)
        if [ "$file" = "$valid_file" ]
            set VALID 1
        end
    end
    if [ "$VALID" = "0" ]
        echo "Deleting $file"
        rm -r $file
    end
end

# Copy the new binaries
echo Copying newly built binaries
cp -rv ../../ui/dist/* .

# Embed Dependencies
cd ../..
./embed_resources.fish

# Update Deps
./update_deps.fish

# Build
./build.fish