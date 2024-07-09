import { createFileRoute,  getRouteApi,  useRouter,  } from '@tanstack/react-router'

import { FilePathListItem, useFilesList } from '../api/fetch_files_list'
import { useEffect, useState } from 'react'
type FileSearch = {
  name ?:string
}
export const Route = createFileRoute('/')({
  component: Index,
  validateSearch: (search: Record<string, unknown>): FileSearch => {
    // validate and parse the search params into a typed state
    return {
      name: (search.name as string),
    }
  },
})


function Index() {
  const router = useRouter()
    const searchParams = Route.useSearch()
  
    const {data,isLoading} = useFilesList(searchParams)
    const [search_term, setSearchTerm] = useState(searchParams.name)
    const [files_list, setFilesList] = useState<FilePathListItem[]>()
   
    useEffect(()=>{
      if(data?.data){
        console.log(`files_list fetched`)
        setFilesList(data.data)
        if(search_term !== searchParams.name) {
          if(search_term){
            router.navigate({
              search:{name:search_term},
              to:Route.fullPath,
              state:true,

            })
          }else{
            router.navigate({
              to:Route.fullPath,
              state:true,
            })
          }
        }

      }
    },[data,search_term])

    

    

  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
      <div className='w-full flex justify-items-center'>
          <form onSubmit={(v)=>{
            v.preventDefault()

            const {name:{value}} = v?.target as EventTarget & {name:{value?:string}}
            
            // console.log({name_v,v})
            setSearchTerm(value)
          }}>
            <input type='text' name="name" className='min-w-[8rem] p-2 border-blue-100 border-2 rounded-lg' defaultValue={searchParams?.name||''}
              onChange={(v)=>{
                v.preventDefault()
                const val = v.target.value as string|undefined
                setSearchTerm(val)
              }}
            />
            <button type='submit' />
          </form>
        </div>
      <div className='flex items-center w-full'>
       
      <div>
        {isLoading && <div> fetching files list from database... </div> }
      {
        files_list?.slice(100).map((v,i)=>(<div key={`path_item_${i}`} className='text-wrap'> {v.path} </div>))
      }
      </div>
      </div>
    </div>
  )
}

