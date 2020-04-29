import axios from 'axios'

var baseURL = "/" || "http://127.0.0.1:8088";

const request = axios.create({
    baseURL: baseURL,
    withCredentials: false,
    timeout: 5000
})
const pageSize = 20

export function fetchProjects(page = 1, page_size = pageSize) {
    return request.get("/api/projects", {
        params: {
            page: page,
            page_size: page_size
        }
    })
}

export function fetchProject(id) {
    return request.get(`/api/project/${id}`)
}

export function startProject(id) {
    return request.get(`/api/start/${id}`)
}

export function stopProject(id) {
    return request.get(`/api/stop/${id}`)
}

export function fetchVuls(pid, page = 1, page_size = pageSize) {
    var params = {
        page: page,
        page_size: page_size
    }
    if (/^[0-9]+.?[0-9]*$/.test(pid)) {
        params["id"] = pid
    } else {
        params["domain"] = pid
    }
    return request.get(`/api/vuls`, {
        params: params
    })
}

export function createProject(form) {
    return request.post("/api/project", {
        name: form.name,
        domain: form.domain,
        plugins: form.plugins.toString(),
    })
}