package curseForge

import (
    "encoding/json"
    "net/http"
    "os"
    "io"

    // "fmt"
)

type SearchMod_Author__struct struct {
    Id int 
    Name string
    Username string
    IsEarlyAccessAuthor bool
}
type SearchMod_Categories__struct struct {
    Id int
    Datemodified string
    Gameid int
    Iconurl string
    Name string
    Slug string
    Url string
    Classid int
    Isclass bool
    Parentcategoryid int
}
type SearchMod_class__struct struct {
    Id int
    DateModified string
    GameId int
    IconUrl string
    Name string
    Slug string
    Url string
    ClassId int
    DisplayIndex int
    IsClass bool
    ParentCategoryId int
}
type SearchMod_LatestFileDetails__struct struct {
    Id int
    GameVersions []string
    GameVersionTypeIds []int
}
type SearchMod_GameVersion__struct struct {
    Id int
    name string
}
type SearchMod_Files__struct struct {
    FileName string
    Id int
    DateCreated string
    DateModified string
    DisplayName string
    ReleaseType int
    GameVersions []string
    EarlyAccessEndDate string
    GameVersionTypeIds []int
    IsEarlyAccessContent bool
}
type SearchMod_WebsiteRecentFiles__struct struct {
    GameVersion SearchMod_GameVersion__struct
    Files []SearchMod_Files__struct
}
type SearchMod_main__struct struct {
    Id int
    Author SearchMod_Author__struct
    AvatarUrl string
    Categories []SearchMod_Categories__struct
    Class SearchMod_class__struct
    CreationDate int
    Downloads int
    GameVersion string
    Name string
    Slug string
    Summary string
    UpdateDate int
    ReleaseDate int
    FileSize int
    IsClientCompatible bool
    LatestFileDetails SearchMod_LatestFileDetails__struct
    HasEarlyAccessFiles bool
    HasLocalization bool
    Status int
    WebsiteRecentFiles []SearchMod_WebsiteRecentFiles__struct
}
type SearchMod_root__struct struct {
    Data []SearchMod_main__struct
}
func SearchMod(searching string) (*SearchMod_root__struct,error) {
    r,err := http.Get("https://www.curseforge.com/api/v1/mods/search?gameId=432&index=0&filterText="+searching+"&pageSize=20&sortField=1")
    if err != nil {
        return nil,err
    }
    body,err := io.ReadAll(r.Body)
    if err != nil {
        return nil,err
    }
    j := SearchMod_root__struct{}
    json.Unmarshal([]byte(string(body)),&j)
    return &j,nil
}

type GetFiles_user__struct struct {
    Username string
    Id int
    TwitchAvatarUrl string
    DisplayName string
}

type GetFiles_files__struct struct {
    Id int
    DateCreated string
    DateModified string
    FileLength int
    FileName string
    Status int
    GameVersions []string
    GameVersionTypeIds []int
    ReleaseType int
    TotalDownloads int
    User GetFiles_user__struct
    AdditionalFilesCount int
    HasServerPack bool
    AdditionalServerPackFilesCount int
    IsEarlyAccessContent bool
}

type GetFiles_root__struct struct {
    Data []GetFiles_files__struct
    Pagination map[string]int
}

func GetFiles(modId string) (*GetFiles_root__struct,error) {
    r,_ := http.Get("https://www.curseforge.com/api/v1/mods/"+modId+"/files?pageIndex=0&pageSize=999999&sort=dateCreated&sortDescending=true&removeAlphas=true")
    body,_ := io.ReadAll(r.Body)
    j := GetFiles_root__struct{}
    json.Unmarshal(body,&j)
    return &j,nil
}

func Download(modId,downloadId,filename, path string) (bool,error) {
    r,err := http.Get("https://www.curseforge.com/api/v1/mods/"+modId+"/files/"+downloadId+"/download")
    if err != nil {
        return false,err
    }
    body,err := io.ReadAll(r.Body)
    if err != nil {
        return false,err
    }
    err = os.WriteFile(path+filename, body, 0644)
    if err != nil {
        return false,err
    }
    return true,nil
}