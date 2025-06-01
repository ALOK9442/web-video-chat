import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

export default function DashBoard() {
  const [InputValue, setInputValue] = useState<string>("");
  const navigate = useNavigate();
  const localVideoRef = useRef<HTMLVideoElement | null>(null);
  const handleNavigate = () => {
    if (InputValue && InputValue.trim() != "") {
      navigate(`/lobby/${encodeURIComponent(InputValue.trim())}`);
    }
  };

  useEffect(() => {
    navigator.mediaDevices
      .getUserMedia({ video: true, audio: true })
      .then((stream) => {
        if (localVideoRef.current) {
          localVideoRef.current.srcObject = stream;
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);

  return (
    <div>
      <video
        autoPlay
        playsInline
        muted
        ref={localVideoRef}
        style={{
          width: "50%",
          borderColor: "orange",
          borderWidth: "10px",
          borderStyle: "solid",
        }}
      />
      <input
        type="text"
        value={InputValue}
        placeholder="Your Name"
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            handleNavigate();
          }
        }}
      />
      <button onClick={handleNavigate}>Submit</button>
    </div>
  );
}
