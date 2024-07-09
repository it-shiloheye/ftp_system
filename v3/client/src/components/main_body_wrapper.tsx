import React, { ReactNode } from "react"
import { Link } from '@tanstack/react-router'
export function NavigationBar({ children }: { children?: ReactNode }) {

    const UserLinks: { address: string, label: string }[] = [{
        address: "/", label: "Home"
    }, {
        address: "/about", label: "About"
    }]


    if (children) {
        return (<div className="w-full flex p-2">
           <Title/>
            {UserLinks.map((v,i) => (<NavLink  key={`nav_link_${i}`}  address={v.address} label={v.label} />))}

            {children}

        </div>)
    }
    return (<div className="w-full flex p-2">
       <Title/>
        {UserLinks.map((v,i) => (<NavLink  key={`nav_link_${i}`}  address={v.address} label={v.label} />))}
    </div>)
}


export default function MainBody({ children }: { children: ReactNode }) {


    return (<div className="h-screen w-screen overflow-x-hidden">
        {children}
    </div>)
}


function Title(){
    return (
        <h1 className="items-center mr-auto text-[2rem]">
        Shiloh File Handler
    </h1>
    )
}

function NavLink({address, label}:{address:string,label:string}){

    return (<Link   to={address} className="[&.active]:font-bold items-center capitalize mx-2 h-full">
        {label}
    </Link>)
}