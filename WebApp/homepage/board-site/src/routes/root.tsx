import boardshapesLogo from "../img/boardshapes.png";
import "./root.css";
import { Outlet, Link } from "react-router-dom";

function Root() {
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
