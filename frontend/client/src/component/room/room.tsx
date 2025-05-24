import { useParams } from "react-router-dom";

export default function Room() {
  const param = useParams();
  console.log(param);
  return <div>room</div>;
}
