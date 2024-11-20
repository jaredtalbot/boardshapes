import { Link } from "react-router-dom";
import "./root.css";

function About() {
  return (
    <>
      <div className="body1">
        <h1>About The Devs:</h1>
        <p className="AboutTitle">Jared</p>
        <p className="AboutDescription">
          Jared is a Fifth-Semester Junior majoring in Computer Science. He is
          the project leader!
        </p>
        <p className="AboutTitle">Cohen</p>
        <p className="AboutDescription">
          Cohen is a Fifth-Semester Junior majoring in Computer Science. He made
          this website!
        </p>
        <p className="AboutTitle">Luke</p>
        <p className="AboutDescription">
          Luke is a Fifth-Semester Junior majoring in Computer Science. He's
          hard at work on the CLI tool!
        </p>
        <p className="AboutTitle">Zach</p>
        <p className="AboutDescription">
          Zach is a Fifth-Semester Junior majoring in Computer Science. He made
          the Boardwalk mobile port!
        </p>
        <Link to={`/`} className="link">
          Return Home!
        </Link>
      </div>
    </>
  );
}

export default About;
