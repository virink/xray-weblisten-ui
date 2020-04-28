import request from './request'

const pageSize = 20

export function fetchProjects(page = 1, page_size = pageSize) {
    return request.get("/projects", {
        params: {
            page: page,
            page_size: page_size
        }
    })
}

export function fetchProject(id) {
    return request.get(`/project/${id}`)
}

export function startProject(id) {
    return request.get(`/start/${id}`)
}

export function stopProject(id) {
    return request.get(`/stop/${id}`)
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
    return request.get(`/projects`, {
        params: params
    })
}