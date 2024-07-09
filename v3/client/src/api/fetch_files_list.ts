import axios from "axios";
import { useQuery } from "@tanstack/react-query";

const files_list_route = `/api/files_list`;

const getFilesList = async ({name}:{name?:string}) => {
	const { data } = await axios.get(`${files_list_route}${name?"?name="+encodeURI(name):""}`);
	return data as {
		success: boolean;
		data: FilePathListItem[];
	};
};

export function useFilesList(search_term:{name?:string}) {
	return useQuery({
		queryKey: ["files_list",search_term],
		queryFn: () => getFilesList(search_term),
	});
}

export type FilePathListItem = {
	hash: string;
	mod_time: Date;
	path: string;
	type: string;
};
