import { Link } from "react-router-dom";
import "./root.css";
import { useMemo } from "react";
import pirateHat from "../img/pirate.png";
import boardgoggle from "../img/boardgoggle.png";
import jester from "../img/jester.png";
import wegahat from "../img/wegahat.png";
import colander from "../img/colander.png";

function About() {
  const todayHash = useMemo(() => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return today
      .toISOString()
      .split("")
      .reduce((a, b) => {
        a = (a << 5) - a + b.charCodeAt(0);
        return a & a;
      }, 0);
  }, []);

  return (
    <>
      <div className="body1">
        <h1>About The Devs:</h1>
        <div className="AboutGroup">
          <a className="hatlink" href={`/boardwalk?unlock=${todayHash}`}>
            <img className="hat" src={jester}></img>
          </a>
          <div className="AboutTitleAndDescription">
            <p className="AboutTitle">Jared</p>
            <p className="AboutDescription">
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Non
              veritatis dolorum assumenda officia nobis modi eos, neque enim
              sapiente laborum voluptates asperiores quis reiciendis maxime
              ducimus mollitia ullam? Placeat, tenetur. Lorem ipsum dolor sit,
              amet consectetur adipisicing elit. Assumenda est accusamus
              temporibus molestiae in nobis illo officia officiis obcaecati
              tenetur sequi itaque mollitia ut placeat facere, similique
              laudantium doloribus velit.
            </p>
          </div>
        </div>
        <div className="AboutGroup">
          <img className="hat" src={boardgoggle}></img>
          <div className="AboutTitleAndDescription">
            <p className="AboutTitle">Cohen</p>
            <p className="AboutDescription">
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Non
              veritatis dolorum assumenda officia nobis modi eos, neque enim
              sapiente laborum voluptates asperiores quis reiciendis maxime
              ducimus mollitia ullam? Placeat, tenetur. Lorem ipsum dolor sit,
              amet consectetur adipisicing elit. Assumenda est accusamus
              temporibus molestiae in nobis illo officia officiis obcaecati
              tenetur sequi itaque mollitia ut placeat facere, similique
              laudantium doloribus velit.
            </p>
          </div>
        </div>
        <div className="AboutGroup">
          <img className="hat" src={wegahat}></img>
          <div className="AboutTitleAndDescription">
            <p className="AboutTitle">Luke</p>
            <p className="AboutDescription">
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Non
              veritatis dolorum assumenda officia nobis modi eos, neque enim
              sapiente laborum voluptates asperiores quis reiciendis maxime
              ducimus mollitia ullam? Placeat, tenetur. Lorem ipsum dolor sit,
              amet consectetur adipisicing elit. Assumenda est accusamus
              temporibus molestiae in nobis illo officia officiis obcaecati
              tenetur sequi itaque mollitia ut placeat facere, similique
              laudantium doloribus velit.
            </p>
          </div>
        </div>
        <div className="AboutGroup">
          <img className="hat" src={pirateHat}></img>
          <div className="AboutTitleAndDescription">
            <p className="AboutTitle">Zachary</p>
            <p className="AboutDescription">
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Non
              veritatis dolorum assumenda officia nobis modi eos, neque enim
              sapiente laborum voluptates asperiores quis reiciendis maxime
              ducimus mollitia ullam? Placeat, tenetur. Lorem ipsum dolor sit,
              amet consectetur adipisicing elit. Assumenda est accusamus
              temporibus molestiae in nobis illo officia officiis obcaecati
              tenetur sequi itaque mollitia ut placeat facere, similique
              laudantium doloribus velit.
            </p>
          </div>
        </div>
        <div className="AboutGroup">
          <img className="hat" src={colander}></img>
          <div className="AboutTitleAndDescription">
            <p className="AboutTitle">Jonathan</p>
            <p className="AboutDescription">
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Non
              veritatis dolorum assumenda officia nobis modi eos, neque enim
              sapiente laborum voluptates asperiores quis reiciendis maxime
              ducimus mollitia ullam? Placeat, tenetur. Lorem ipsum dolor sit,
              amet consectetur adipisicing elit. Assumenda est accusamus
              temporibus molestiae in nobis illo officia officiis obcaecati
              tenetur sequi itaque mollitia ut placeat facere, similique
              laudantium doloribus velit.
            </p>
          </div>
        </div>
        <br></br>
        <br></br>
        <Link to={`/`} className="link">
          Return Home!
        </Link>
      </div>
    </>
  );
}

export default About;
