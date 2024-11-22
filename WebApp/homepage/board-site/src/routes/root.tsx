import boardmeshLogo from "../img/boardmesh.svg";
import "./root.css";
import { Outlet, Link } from "react-router-dom";

function Root() {
  return (
    <>
      <div className="body1">
        <img src={boardmeshLogo} className="logoMain" alt="Boardwalk logo" />
        <p>
          Boardmesh is a free and open-source API that can be used to create
          suitable, color-corrected, collidable shapes for physics simulators
          and simplify an image into a restricted color palette consisting of
          black, white, red, green, and blue.
        </p>
        <p>
          Tapping "Boardwalk" will take you to our game implementation of the
          Boardmesh API, Boardwalk. Boardwalk is a single-screen platformer that
          allows the player to traverse their newly created meshes.
        </p>
        <p>
          Tapping "Boardmesh" will take you to the User Manual for Boardwalk.
        </p>
        <p>
          Tapping the third logo will take you to our Physics Sim, Boardbox.
          This allows you to interact with your new meshes in a very direct
          manner.
        </p>
        <p>Tapping the Github logo will take you to our GitHub repository.</p>
        <p className="read-the-docs">
          If you wanna know more about the developers, click the link below!
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
