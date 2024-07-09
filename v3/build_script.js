const fs = require("fs");

const cp = require("child_process");

const build_promise = () => {
	return new Promise((resolve, reject) => {
        console.log("building executable: ")
		cp.exec(
			"powershell.exe scripts/build_script.ps1",
			(error, stdout, stderr) => {
				if (error || stderr) {
					if (error) {
						console.log(error);
					}
					if (stderr) {
						console.log(stderr);
					}
					reject();
					return;
				}

                console.log(stdout)
				console.log("success");
                resolve()
			},
		);
	});
};

const date = (() => {
	const d = Date();
    const d_split = d.split(" ")
    const [_, month,day,year,time] = d_split
    const time_clean = time.split(":").join("")
    const date_full = `${year}_${month}_${day}_${time_clean}`
	console.log(date_full);
    return date_full
    
})();

console.log("hello world");

const count = fs.readdirSync("./build",{
    withFileTypes:true,
    recursive:true
}).filter(v=>v.isFile() && v.name.includes("ftp_server_build")).length

const rename_promise = () => {
	return new Promise((resolve, reject) => {
        console.log("renaming executable: ")
		cp.exec(
			`mv build/ftp_server_build.exe build/ftp_server_build_${count.toString().padStart(5,"0")}_${date}.exe`,
			(error, stdout, stderr) => {
				if (error || stderr) {
					if (error) {
						console.log(error);
					}
					if (stderr) {
						console.log(stderr);
					}
					reject();
					return;
				}

                console.log(stdout)
				console.log("success");
                resolve()
			},
		);
	});
};

build_promise().then(()=>rename_promise()).then(()=>{
    console.log("build successful")
}).catch((err)=>{})