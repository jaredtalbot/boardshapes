import boardwalkLogo from "../img/boardwalk.svg";
import boardmeshLogo from "../img/boardmesh.svg";
import githubLogo from "../img/github.svg";
import boardboxLogo from "../img/BoardBox.svg";
import "./root.css";
import { Outlet, Link } from "react-router-dom";
import Headroom from "react-headroom";

function Root() {
  return (
    <>
      <Headroom
        wrapperStyle={{ marginBottom: 35 }}
        style={{
          background: "rgb(220, 220, 220)",
        }}
      >
        <div style={{ padding: 10 }}>
          <h1
            style={{
              margin: 0,
              color: "rgb(252, 253, 254)",
            }}
          >
            <div>
              <a href="/boardwalk/" target="_blank">
                <img
                  src={boardwalkLogo}
                  className="logo"
                  alt="Boardwalk logo"
                />
              </a>
              <a href="/manual/" target="_blank">
                <img
                  src={boardmeshLogo}
                  className="logo"
                  alt="Boardmesh logo"
                />
              </a>
              <a href="/boardbox/" target="_blank">
                <img src={boardboxLogo} className="logo" alt="Boardbox logo" />
              </a>
              <a
                href="https://github.com/codeJester27/cmps401fa2024"
                target="_blank"
              >
                <img src={githubLogo} className="logo" alt="Github logo" />
              </a>
            </div>
          </h1>
        </div>
      </Headroom>
      <div className="body1">
        <img src={boardmeshLogo} className="logoMain" alt="Boardwalk logo" />
        <h1>
          Boardmesh is a free and open-source API that can be used to create
          suitable, color-corrected, collidable shapes for physics simulators
          and simplify an image into a restricted color palette consisting of
          black, white, red, green, and blue.
        </h1>
        <h1>
          Tapping "Boardwalk" will take you to our game implementation of the
          Boardmesh API, Boardwalk. Boardwalk is a single-screen platformer that
          allows the player to traverse their newly created meshes.
        </h1>
        <h1>
          Tapping "Boardmesh" will take you to the User Manual for Boardwalk.
        </h1>
        <h1>
          Tapping the third logo will take you to our Physics Sim, Boardbox.
          This allows you to interact with your new meshes in a very direct
          manner.
        </h1>
        <h1>Tapping the Github logo will take you to our GitHub repository.</h1>
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
