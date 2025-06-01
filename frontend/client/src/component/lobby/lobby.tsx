import { useParams } from "react-router-dom";
import "./lobby.css";
import { useEffect, useRef, useState } from "react";

export default function Lobby() {
  // const url = "htt"
  const params = useParams();
  const name = params.name;
  const peerRef = useRef<RTCPeerConnection | null>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const localVideoRef = useRef<HTMLVideoElement | null>(null);
  const remoteVideoref = useRef<HTMLVideoElement | null>(null);
  const newIceCandidateQueue = useRef<RTCIceCandidateInit[]>([]);
  const remoteDescriptionStatus = useRef(false);

  useEffect(() => {
    // console.log(ws)
    const peer = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
    });
    console.log(peer);

    peerRef.current = peer;

    console.log("Initial connection state:", peer.connectionState);

    navigator.mediaDevices
      .getUserMedia({ video: true, audio: true })
      .then((stream) => {
        console.log("stream", stream);
        if (localVideoRef.current) {
          console.log("inside current", localVideoRef);
          localVideoRef.current.srcObject = stream;
        } else {
          console.log("else");
        }
        // stream.getTracks().forEach((track) => {
        //   console.log("tracks", track);
        //   peer.addTrack(track, stream);
        // });
      })
      .catch((err) => {
        console.error("Error getting user media:", err);
      });

    peer.ontrack = async (e) => {
      console.log(e);
      console.log(e.streams);
      if (remoteVideoref.current) {
        console.log("hehe", remoteVideoref);
        remoteVideoref.current.srcObject = e.streams[0];
      }
    };
    const ws: WebSocket | null = new WebSocket("ws://localhost:8080/ws");
    setTimeout(() => {
      const receivers = peer.getReceivers();
      receivers.forEach((r) => {
        console.log("Receiver track:", r.track);
      });
    }, 2000);

    peer.onicecandidate = (e) => {
      if (e.candidate !== null && ws?.readyState === WebSocket.OPEN) {
        console.log(e.candidate);
        ws.send(
          JSON.stringify({
            type: "ice-candidate",
            candidate: e.candidate,
          })
        );
      }
    };

    peer.onconnectionstatechange = () => {
      console.log("Connection state changed:", peer.connectionState);
    };

    peer.oniceconnectionstatechange = () => {
      console.log("ICE connection state changed:", peer.iceConnectionState);
    };

    peer.onsignalingstatechange = () => {
      console.log("Signaling state changed:", peer.signalingState);
    };

    ws.onopen = async () => {
      console.log("WebSocket connection opened");
      setSocket(ws);
    };

    ws.onmessage = async (e) => {
      console.log(e);
      // console.log(e.data);

      const outerData = JSON.parse(e.data);
      console.log(outerData);

      // const data =
      //   outerData.type === "offer" ? outerData.offer : outerData.data;
      // console.log(data);
      switch (outerData.type) {
        case "waiting":
          console.log("waiting");
          break;
        case "caller": {
          const stream = await navigator.mediaDevices.getUserMedia({audio:true, video:true})
          if (localVideoRef.current){
            localVideoRef.current.srcObject = stream
          }
          stream.getTracks().forEach(track => {
            peer.addTrack(track,stream)
          });
          const ClientOffer = await peer.createOffer();
          console.log(ClientOffer);
          await peer.setLocalDescription(ClientOffer);
          ws.send(
            JSON.stringify({
              type: "offer",
              ClientOffer,
            })
          );
          break;
        }
        case "offer":
          console.log("here", outerData);
          if (outerData) {
            await peer.setRemoteDescription(
              new RTCSessionDescription(outerData.ClientOffer)
            );
            const stream = await navigator.mediaDevices.getUserMedia({audio:true, video:true})
            if(localVideoRef.current){
              localVideoRef.current.srcObject = stream
            }
            stream.getTracks().forEach(track => {
              peer.addTrack(track,stream)
            });
            remoteDescriptionStatus.current = true;

            newIceCandidateQueue.current.forEach((element) => {
              peer.addIceCandidate(new RTCIceCandidate(element));
            });

            console.log(outerData);
            const answer = await peer.createAnswer();
            console.log(answer);
            await peer.setLocalDescription(answer);
            ws.send(
              JSON.stringify({
                type: "answer",
                answer,
              })
            );
          }
          break;
        case "answer":
          // if (outerData) {
          console.log(outerData);
          await peer.setRemoteDescription(
            new RTCSessionDescription(outerData.answer)
          );
          remoteDescriptionStatus.current = true;
          newIceCandidateQueue.current.forEach((element) => {
            peer.addIceCandidate(new RTCIceCandidate(element));
          });
          // }
          break;
        case "skip":
          peerRef.current?.close();
          peerRef.current = null;
          break;
        case "ice-candidate":
          // if (outerData) {
          if (!remoteDescriptionStatus.current) {
            newIceCandidateQueue.current.push(outerData.candidate);
          } else {
            await peer.addIceCandidate(
              new RTCIceCandidate(outerData.candidate)
            );
          }
          // }
          break;
        case "chat":
          console.log("chat:", outerData);
          break;
        default:
          break;
      }
    };
    ws.onclose = (event) => {
      console.warn("WebSocket closed:", event, event.code, event.reason);
    };
    return () => {
      ws.close();
      peer.close();
    };
  }, []);
  console.log(params);
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <video
        ref={localVideoRef}
        style={{
          width: "50%",
          borderColor: "orange",
          borderWidth: "10px",
          borderStyle: "solid",
        }}
        autoPlay
        playsInline
        muted
      />
      <video
        ref={remoteVideoref}
        autoPlay
        playsInline
        style={{
          width: "50%",
          borderColor: "orange",
          borderWidth: "10px",
          borderStyle: "solid",
        }}
      />
    </div>
  );
}
