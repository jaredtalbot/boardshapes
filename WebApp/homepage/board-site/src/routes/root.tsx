import boardshapesLogo from "../img/boardshapes.png";
import "./root.css";
import { Outlet, Link } from "react-router-dom";
import { useState, useRef, useEffect } from "react";

function Root() {
  const [regionData, setRegionData] = useState<RegionData[] | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  interface RegionData {
    regionNumber: number;
    regionColor: { R: number; G: number; B: number; A: number };
    regionColorString: string;
    cornerX: number;
    cornerY: number;
    regionImage: string;
    mesh: { x: number; y: number }[];
  }

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
    setRegionData(null);

    const formData = new FormData();
    formData.append("image", file);

    try {
      const response = await fetch("/api/build-level", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("Server response:", response.status, errorText);
        throw new Error(
          `Failed to process image: ${response.status} - ${errorText}`
        );
      }

      const data: RegionData[] = await response.json();
      setRegionData(data);
    } catch (error) {
      console.error("Error details:", error);
      setErrorMessage(
        error instanceof Error ? error.message : "An unknown error occurred."
      );
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (!regionData || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let maxX = 0;
    let maxY = 0;
    regionData.forEach((region) => {
      region.mesh.forEach((vertex) => {
        maxX = Math.max(maxX, vertex.x + region.cornerX);
        maxY = Math.max(maxY, vertex.y + region.cornerY);
      });
    });

    canvas.width = maxX + 20;
    canvas.height = maxY + 20;

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    regionData.forEach((region) => {
      const img = new Image();
      img.src = `data:image/png;base64,${region.regionImage}`;
      img.onload = () => {
        ctx.drawImage(img, region.cornerX, region.cornerY);

        ctx.beginPath();
        const vertices = region.mesh;
        if (vertices.length > 0) {
          ctx.moveTo(
            vertices[0].x + region.cornerX,
            vertices[0].y + region.cornerY
          );
          for (let i = 1; i < vertices.length; i++) {
            ctx.lineTo(
              vertices[i].x + region.cornerX,
              vertices[i].y + region.cornerY
            );
          }
          ctx.closePath();
          ctx.fillStyle = "rgba(112, 112, 255, 0.46)";
          ctx.fill();
          ctx.strokeStyle = "rgb(112, 112, 255)";
          ctx.lineWidth = 2;
          ctx.stroke();
        }
      };
    });
  }, [regionData]);

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
          <h3>Collision Shapes Visualization</h3>
          <input
            type="file"
            accept="image/png,image/jpeg"
            onChange={handleImageUpload}
            disabled={isLoading}
          />
          {isLoading && <p className="loading">Processing image...</p>}
          {errorMessage && <p className="error">{errorMessage}</p>}
          {regionData && (
            <div className="simplified-image-container">
              <h4>Collision Shapes:</h4>
              <canvas ref={canvasRef} className="collision-canvas" />
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
