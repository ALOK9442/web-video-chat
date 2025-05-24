import { useState } from "react";
import { useNavigate } from "react-router-dom"

export default function DashBoard() {
    const [InputValue, setInputValue] = useState<string>("")
    const navigate = useNavigate();
    const handleNavigate = () =>{
        if(InputValue && InputValue.trim() != "") {
            navigate(`/lobby/${encodeURIComponent(InputValue.trim())}`)
        }
    }
    return (
        <div>
            <input type="text" value={InputValue} placeholder="Your Name" onChange={(e)=>setInputValue(e.target.value)} />
            <button
            onClick={handleNavigate}
            >Submit</button>
        </div>
    )
}