import boardshapesLogo from "../img/boardshapes.png";
import "./root.css";
import { Outlet, Link } from "react-router-dom";
import { useState } from "react";

function Root() {
  const [simplifiedImage, setSimplifiedImage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleImageUpload = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) {
      setErrorMessage("No file selected.");
      return;
    }

    if (!["image/png", "image/jpeg"].includes(file.type)) {
      setErrorMessage("Please upload a PNG or JPEG image.");
      return;
    }

    setIsLoading(true);
    setErrorMessage(null);

    const formData = new FormData();
    formData.append("image", file);

    try {
      const response = await fetch("/api/simplify", {
        method: "POST", //kid named follows the manual
        body: formData,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("Server response:", response.status, errorText);
        throw new Error(
          `Failed to simplify image: ${response.status} - ${errorText}`
        );
      }

      const imageBlob = await response.blob();
      const imageUrl = URL.createObjectURL(imageBlob);
      setSimplifiedImage(imageUrl);
    } catch (error) {
      console.error("Error details:", error);
      setErrorMessage(
        error instanceof Error ? error.message : "An unknown error occurred."
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <div className="body1">
        <img src={boardshapesLogo} className="logoMain" alt="Boardwalk logo" />
        <div className="boardshapes-text-container">
          <p id="paraJustify">
            <strong>Boardshapes</strong> (formerly "Boardmesh") is a free and
            open-source API that can be used to create suitable,
            color-corrected, collidable shapes for physics simulators and
            simplify an image into a restricted color palette consisting of
            black, white, red, green, and blue.
          </p>
          <ul id="bodylist">
            <li>
              Tapping <strong>"Boardwalk"</strong> will take you to our game
              implementation of the Boardshapes API, Boardwalk. Boardwalk is a
              single-screen platformer that allows the player to traverse their
              newly created shapes.
            </li>
            <li>
              Tapping <strong>"Boardshapes"</strong> will take you to the User
              Manual for Boardwalk.
            </li>
            <li>
              Tapping the third logo will take you to our Physics Sim,{" "}
              <strong>Boardbox</strong>. This allows you to interact with your
              new shapes in a very direct manner.
            </li>
            <li>
              Tapping the <strong>Github</strong> logo will take you to our
              GitHub repository.
            </li>
          </ul>
        </div>

        <div className="image-upload-section">
          <h3>Try Image Simplification</h3>
          <input
            type="file"
            accept="image/png,image/jpeg"
            onChange={handleImageUpload}
            disabled={isLoading}
          />
          {isLoading && <p>Processing image...</p>}
          {errorMessage && (
            <p style={{ color: "red", marginTop: "10px" }}>{errorMessage}</p>
          )}
          {simplifiedImage && (
            <div className="simplified-image-container">
              <h4>Simplified Image:</h4>
              <img
                src={simplifiedImage}
                alt="Simplified version"
                style={{ maxWidth: "100%", marginTop: "10px" }}
              />
            </div>
          )}
        </div>

        <p className="read-the-docs">
          <strong>
            If you want to know more about the developers, click the link below!
          </strong>
        </p>
        <Link to={`/about`} className="link">
          Click Here!
        </Link>
      </div>
      <div id="detail">
        <Outlet />
      </div>
    </>
  );
}

export default Root;
